import {App} from "@/App.tsx";
import {StrictMode} from "react";
import {createRoot} from "react-dom/client";
import "@/util/i18n.ts";

createRoot(document.getElementById("root")!).render(
    <StrictMode>
        <App/>
    </StrictMode>,
);
