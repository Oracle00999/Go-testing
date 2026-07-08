# Portiq Go Test API

This project exists only to test Portiq, the deployment platform being built locally by me.

It is a small Go HTTP API used to verify Go auto-detection, generated Dockerfiles, required environment variables, health checks, runtime logs, and redeploys.

It intentionally has no Dockerfile so Portiq should detect Go from `go.mod` and generate one.

## Required Environment Variables

- `APP_NAME`: display name for the API
- `API_KEY`: required as the `x-api-key` header for protected routes
- `GO_ENV`: runtime environment label
- `RELEASE`: release label returned by the API
- `PORT`: container port, use `3000`

## Deploy Settings

- Internal port: `3000`
- Build command: leave empty
- Start command: leave empty
- Dockerfile path: leave empty
