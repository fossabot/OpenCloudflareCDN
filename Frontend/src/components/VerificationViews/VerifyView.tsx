import {LoadingSpinner} from "@/components/LoadingSpinner/LoadingSpinner";
import {Turnstile} from "@marsidev/react-turnstile";
import type {TFunction} from "i18next";
import i18n from "i18next";
import React from "react";

interface VerifyViewProps {
    t: TFunction;
    domain: string;
    isTurnstileLoaded: boolean;
    noCaptcha: boolean;
    siteKey: string | null | undefined;
    onSuccess: (token: string) => void;
    onError: () => void;
    onExpire: () => void;
    onLoad: () => void;
}

const VerifyView: React.FC<VerifyViewProps> = ({
                                                   t,
                                                   domain,
                                                   isTurnstileLoaded,
                                                   noCaptcha,
                                                   siteKey,
                                                   onSuccess,
                                                   onError,
                                                   onExpire,
                                                   onLoad,
                                               }) => (
    <div>
        <p className="h2 spacer-bottom">{t('verify')}</p>
        <div className="turnstile-container">
            {(!isTurnstileLoaded || noCaptcha) && <LoadingSpinner/>}
            {!noCaptcha && (
                <Turnstile
                    siteKey={siteKey || '1x00000000000000000000AA'}
                    onSuccess={onSuccess}
                    onError={onError}
                    onExpire={onExpire}
                    onLoad={onLoad}
                    options={{
                        language: i18n.language,
                    }}
                />
            )}
        </div>
        <p className="core-msg spacer-top">{t('check_connection', {domain})}</p>
    </div>
);

export default React.memo(VerifyView);