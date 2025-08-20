export interface Response<T> {
    msg: string;
    data?: T;
    raw?: globalThis.Response;
}

export interface VerifyPayload {
    token: string;
    rayID: string;
}