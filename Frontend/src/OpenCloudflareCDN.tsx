import ErrorView from "@/components/VerificationViews/ErrorView";
import SuccessView from "@/components/VerificationViews/SuccessView";
import VerifyView from "@/components/VerificationViews/VerifyView";
import {getRootDomain} from "@/util/cloudflare";
import "@/OpenCloudflareCDN.scss";
import i18n from "i18next";
import React, {useEffect, useState} from "react";
import {Trans, useTranslation} from "react-i18next";

export type VerificationStatus = 'verify' | 'success' | 'error' | 'expire';

interface AppProps {
    siteKey: string | null | undefined;
    successCallback: (token: string, rayID: string) => Promise<void>;
    rayID: string | null | undefined;
    status?: VerificationStatus;
}

const OpenCloudflareCDNComponent: React.FC<AppProps> = ({
                                                            rayID,
                                                            siteKey,
                                                            successCallback,
                                                            status: initialStatus = 'verify'
                                                        }) => {
    const {t} = useTranslation();
    const [status, setStatus] = useState<VerificationStatus>(initialStatus);
    useEffect(() => {
        setStatus(initialStatus)
    }, [initialStatus]);
    const [isTurnstileLoaded, setIsTurnstileLoaded] = useState(false);
    const domain = getRootDomain();
    const rid = status === 'error' ? "Unknown" : rayID || "Loading...";
    const noCaptcha = !siteKey;

    document.title = i18n.t('page_title');

    const renderContent = () => {
        if (status === 'success') {
            return <SuccessView t={t} domain={domain}/>;
        }
        if (status === 'error' || status === 'expire') {
            return <ErrorView t={t} status={status}/>;
        }
        return (
            <VerifyView
                t={t}
                domain={domain}
                isTurnstileLoaded={isTurnstileLoaded}
                noCaptcha={noCaptcha}
                siteKey={siteKey}
                onSuccess={(token) => {
                    setStatus('success');
                    void successCallback(token, rid);
                }}
                onError={() => setStatus('error')}
                onExpire={() => setStatus('expire')}
                onLoad={() => setIsTurnstileLoaded(true)}
            />
        );
    };

    return (
        <>
            <div className="main-wrapper" role="main">
                <div className="main-content">
                    <h1 className="zone-name-title h1">{domain}</h1>
                    {renderContent()}
                </div>
            </div>
            <div className="footer" role="contentinfo">
            </div>
            <div className="footer" role="contentinfo">
                <div className="footer-inner">
                    <div className="clearfix diagnostic-wrapper">
                        <div className="ray-id">Ray ID: <code>{rid}</code></div>
                    </div>
                    <div className="text-center" id="footer-text">
                        <Trans i18nKey="provided_by">
                            Performance & security by <a
                            rel="noopener noreferrer"
                            href="https://github.com/Sn0wo2/OpenCloudflareCDN"
                            target="_blank"
                        >OpenCloudflareCDN</a>
                        </Trans>
                    </div>
                </div>
            </div>
        </>
    );
};

export const OpenCloudflareCDN = React.memo(OpenCloudflareCDNComponent);