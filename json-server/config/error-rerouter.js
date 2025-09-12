const getFailRoute = (req) => {
    return req.headers?.cookie?.match(/fail-route=(?<failRoute>[^;]+);/)?.groups
        .failRoute;
};

const getStatusCode = (req) => {
    return req.headers?.cookie?.match(/fail-code=(?<statusCode>[^;]+);/)?.groups
        .statusCode;
};

module.exports = (req, res, next) => {
    if (["POST", "PATCH"].includes(req.method)) {
        const failRoute = getFailRoute(req);

        if (failRoute) {
            req.method = "GET";
            req.url = `/errors/${failRoute}`;
            res.status(getStatusCode(req) ?? 400);
        }
    }
    next();
};
