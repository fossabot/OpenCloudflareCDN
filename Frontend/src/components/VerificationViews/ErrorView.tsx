import type {VerificationStatus} from "@/OpenCloudflareCDN";
import type {TFunction} from "i18next";
import React from "react";

interface ErrorViewProps {
    t: TFunction;
    status: VerificationStatus;
}

const ErrorView: React.FC<ErrorViewProps> = ({t, status}) => (
    <div>
        <div id="challenge-error-text" className="h2">{t(`challenge_${status}`)}</div>
        <div className="spacer"></div>
        <div className="core-msg spacer">{t(`challenge_${status}_description`)}</div>
    </div>
);

export default React.memo(ErrorView);