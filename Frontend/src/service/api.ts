import type {Response as ApiResponse, VerifyPayload} from "@/service/types.ts";

async function safeJsonParse<T>(res: globalThis.Response): Promise<T | null> {
    try {
        return await res.json() as T;
    } catch (err) {
        console.error(err);
        return null;
    }
}

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/ocfc-api/';

export const Verify = async (
    payload: VerifyPayload
): Promise<ApiResponse<{ ec?: string }>> => {
    const res = await fetch(API_BASE_URL + "v1/verify", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(payload),
    });

    const parsed = await safeJsonParse<{ msg: string; data: { ec?: string } }>(res);

    if (!parsed) {
        return {
            msg: "failed parse verify response",
            raw: res,
        };
    }

    return {
        msg: parsed.msg,
        data: parsed.data,
        raw: res,
    };
};

export const GetInfo = async (): Promise<ApiResponse<{ siteKey?: string; rayID?: string }>> => {
    const res = await fetch(API_BASE_URL + "v1/info");

    const parsed = await safeJsonParse<{ msg: string; data: { siteKey?: string; rayID?: string } }>(res);

    if (!parsed) {
        return {
            msg: "failed parse info response",
            raw: res,
        };
    }

    return {
        msg: parsed.msg,
        data: parsed.data,
        raw: res,
    };
};
