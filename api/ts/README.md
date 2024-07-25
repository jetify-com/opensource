# Jetify API TypeScript Client

**Warning:** this client is currently in development and subject to change.

## Usage

```bash
# create a new jetify API token
> devbox auth tokens new
```

```typescript
import { ApiClient } from '@jetify/client';

const client = new ApiClient({
  baseUrl: process.env.JETIFY_API_BASE_URL,
  token: process.env.JETIFY_API_TOKEN,
});

const response = await client.sandboxService.createSandbox({
  external_billing_tag: 'my-billing-tag',
  repo: '',
  subdir: '',
  ref: '',
});

console.log(response);

```
