import { useEffect, useState } from "react";
import { OpenCloudflareCDN } from "@/OpenCloudflareCDN.tsx";
import { LoadingSpinner } from "@/components/LoadingSpinner/LoadingSpinner.tsx";
import { genRayID } from "@/util/cloudflare.ts";

export const App = () => {
  const [info, setInfo] = useState<{ siteKey: string; rayID: string } | null>(
    null,
  );

  useEffect(() => {
    const fetchInfo = async () => {
      if (import.meta.env.DEV) {
        setInfo({
          siteKey: import.meta.env.VITE_SITE_KEY,
          rayID: genRayID(),
        });
        return;
      }

      try {
        const response = await fetch("/v1/info");
        const data = await response.json();
        if (data.success) {
          setInfo(data.data);
        } else {
          console.error(data.msg);``
        }
      } catch (err) {
        console.error(err);
      }
    };

    void fetchInfo();
  }, []);

  if (!info) {
    return <LoadingSpinner />;
  }

  return (
    <OpenCloudflareCDN
      siteKey={info.siteKey}
      successCallback={async (token, rayID) => {
        try {
          const response = await fetch("/v1/verify", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              turnstileToken: token,
              rayID: rayID,
            }),
          });
          const data = await response.json();
          if (!data.success) {
            console.error(data.msg);
          }
        } catch (err) {
          console.error(err);
        }
        window.location.reload();
      }}
      rayID={info.rayID}
    />
  );
};
