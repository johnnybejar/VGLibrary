import wretch from "wretch";

interface AccessTokenResponse {
    access_token: string;
    expires_in: number;
    token_type: string;
}

function getAccessToken() {
    const CLIENT_ID = import.meta.env.VITE_TWITCH_CLIENT_ID;
    const CLIENT_SECRET = import.meta.env.VITE_TWITCH_CLIENT_SECRET;

    let response;

    const w = wretch(`https://id.twitch.tv/oauth2/token?client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}&grant_type=client_credentials`)
        .post()
        .forbidden((err) => console.log(err))
        .unauthorized(
            (err, req) => {
                // Handle access token refresh
            })
        .fetchError((err) => console.log(err))
        .res(res => console.log(res))
        .then(res => {
            console.log(res);
        });
    
}

const utils = {
    getAccessToken
};

export default utils;