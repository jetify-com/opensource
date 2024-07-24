# Jetify API TypeScript Client

**Warning:** this client is currently and development and subject to change.

## Usage

```bash
# create a new jetify API token
> devbox auth tokens new
```

```typescript
import { createPromiseClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { SandboxService } from '@jetify/client';

const transport = createConnectTransport({
  baseUrl: process.env.JETIFY_API_URL || 'https://api.jetpack.io',
});

const client = createPromiseClient(SandboxService, transport);

const response = await client.createSandbox(
  {
    external_billing_tag: 'my-billing-tag',
    repo: 'my-repo',
    subdir: 'subdir-in-repo',
    ref: 'branch-or-tag-or-commit',
  },
  {
    headers: {
      Authorization: `Token ${process.env.JETIFY_API_TOKEN}`,
    },
  },
);

console.log(response);
```
