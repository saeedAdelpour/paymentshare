# Third-party browser assets

`internal/httpapi/docs/swagger-ui-bundle.js` and `swagger-ui.css` are unmodified
Swagger UI **v5.34.1** distribution assets from
[swagger-api/swagger-ui](https://github.com/swagger-api/swagger-ui/tree/v5.34.1/dist).
The Apache-2.0 license is retained as `internal/httpapi/docs/LICENSE`.

The app serves these assets locally. When upgrading, download both files and the
license from the same reviewed release, update this notice and the SHA256SUMS
file, then verify Swagger UI can load the contract and execute an operation.
