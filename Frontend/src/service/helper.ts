import type {Response} from "@/service/types.ts";

export const isSuccessResponse = (response: Response<unknown>) => {
    if (!response?.raw?.status) {
        return false;
    }
    const status = response.raw.status
    return status >= 200 && status < 300;
};