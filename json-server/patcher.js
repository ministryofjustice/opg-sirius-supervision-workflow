/**
 * Patcher middleware changes PUT to PATCH, so the object is updated rather than
 * replaced
 */
module.exports = (req, res, next) => {
    if (req.method === "PUT") {
        req.method = "PATCH";
    }
    next();
};
