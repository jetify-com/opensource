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
