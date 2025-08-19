package framework

import (
	"net"
	"net/http"
	"time"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/Sn0wo2/OpenCloudflareCDN/debug"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	middlewarehttp "github.com/Sn0wo2/OpenCloudflareCDN/middleware/http"
	"github.com/gin-gonic/gin"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func Gin() *gin.Engine {
	if !debug.IsDebugging() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	engine.Use(gin.Recovery())

	engine.HandleMethodNotAllowed = true

	return engine
}

func Start(engine *gin.Engine) error {
	g := new(errgroup.Group)

	server := &http.Server{
		Addr:              config.Instance.Server.Address,
		Handler:           engine,
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	if config.Instance.Server.TLS.Cert != "" && config.Instance.Server.TLS.Key != "" {
		g.Go(func() error {
			log.Instance.Info("Listening on TLS", zap.String("address", config.Instance.Server.Address))

			return server.ListenAndServeTLS(config.Instance.Server.TLS.Cert, config.Instance.Server.TLS.Key)
		})
		g.Go(func() error {
			h3Server := &http3.Server{
				Addr:    config.Instance.Server.Address,
				Handler: engine,
				QUICConfig: &quic.Config{
					MaxIdleTimeout:       30 * time.Second,
					HandshakeIdleTimeout: 10 * time.Second,
				},
			}

			log.Instance.Info("Listening on QUIC", zap.String("address", config.Instance.Server.Address))

			return h3Server.ListenAndServeTLS(config.Instance.Server.TLS.Cert, config.Instance.Server.TLS.Key)
		})
	} else {
		log.Instance.Info("Listening on HTTP", zap.String("address", config.Instance.Server.Address))

		g.Go(func() error {
			return server.ListenAndServe()
		})
	}

	if config.Instance.Server.HttpAddress != "" {
		g.Go(func() error {
			redirectEngine := gin.New()
			_, httpsPort, _ := net.SplitHostPort(config.Instance.Server.Address)
			redirectEngine.Use(middlewarehttp.Handler(httpsPort))

			httpServer := &http.Server{
				Addr:              config.Instance.Server.HttpAddress,
				Handler:           redirectEngine,
				ReadHeaderTimeout: 30 * time.Second,
			}

			log.Instance.Info("Listening on HTTP", zap.String("address", config.Instance.Server.HttpAddress))

			return httpServer.ListenAndServe()
		})
	}

	return g.Wait()
}
