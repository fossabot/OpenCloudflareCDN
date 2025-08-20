import type {TFunction} from "i18next";
import React from "react";

interface SuccessViewProps {
    t: TFunction;
    domain: string;
}

const SuccessView: React.FC<SuccessViewProps> = ({t, domain}) => (
    <div>
        <div id="challenge-success-text" className="h2">{t('challenge_success')}</div>
        <div className="spacer"></div>
        <div className="core-msg spacer">{t('waiting', {domain})}</div>
    </div>
);

export default React.memo(SuccessView);