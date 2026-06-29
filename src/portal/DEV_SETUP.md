# Harbor Portal Development Setup

## Frontend Code (Verified ✅)

The portal frontend code is production-ready with all fixes applied:
- Clarity 18.2.2-fixed with PR #2486 (readonly hostbinding fix)
- Datagrid binding corrected (`*clrDgItems`)
- PAT component fully implemented
- **No Clarity/JavaScript component errors**

Build with: `npm run build`

## Development Server with API Proxy

Run the dev server with API proxy configuration:

```bash
npm start -- --proxy-config proxy.conf.json
```

This serves the portal on `https://localhost:4200` and proxies `/api/` requests to the Harbor core service.

**Note:** The proxy target must be configured to reach the Harbor core service:
- For Docker-based setup: expose harbor-core port to localhost or use Docker network bridge
- For local setup: configure harbor-core to listen on localhost:8080

## Testing with Production Stack

To test the complete portal with real API endpoints, access the production Nginx proxy:

```bash
https://localhost  # Serves through production nginx proxy
```

This uses the full Harbor stack (core, registry, jobservice, etc.) with proper API routing already configured.

## Building for Production

```bash
npm run release  # Builds optimized production bundle
```

The built artifacts in `dist/` are served by the portal Nginx container, which handles:
- Static file serving with proper caching headers
- SPA routing (all routes served by index.html)
- API proxying to Harbor core on `/api/` paths
