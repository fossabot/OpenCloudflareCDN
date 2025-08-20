import {OpenCloudflareCDN, type VerificationStatus} from "@/OpenCloudflareCDN.tsx";
import {GetInfo, Verify} from "@/service/api.ts";
import {isSuccessResponse} from "@/service/helper.ts";
import {useCallback, useEffect, useState} from "react";

export const App = () => {
    const [info, setInfo] = useState<{ siteKey?: string; rayID?: string } | null>(null);

    const [status, setStatus] = useState<VerificationStatus>("verify");

    useEffect(() => {
        const fetchInfo = async () => {
            try {
                const resp = await GetInfo();
                if (!isSuccessResponse(resp)) {
                    console.error('Failed to load site info:', resp);
                    setStatus("error")

                    return
                }
                if (!resp.data?.siteKey || !resp.data?.rayID) return
                setInfo({
                    siteKey: resp.data.siteKey,
                    rayID: resp.data.rayID,
                });
            } catch (err) {
                console.error(err);
                setStatus("error")
            }
        };

        void fetchInfo();
    }, []);


    const successCallback = useCallback(async (token: string, rayID: string) => {
        try {
            const resp = await Verify({
                token,
                rayID
            });
            if (!isSuccessResponse(resp)) {
                console.error('Failed to verify:', resp);
                setStatus("error")
            }
        } catch (err) {
            console.error(err);
            setStatus("error")
        }
        window.location.reload();
    }, []);

    return (
        <OpenCloudflareCDN
            siteKey={info?.siteKey}
            successCallback={successCallback}
            rayID={info?.rayID}
            status={status}
        />
    );
};
