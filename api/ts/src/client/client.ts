import {
  createPromiseClient,
  PromiseClient,
  Transport,
} from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { SandboxService } from '../gen/pub/sandbox/v1alpha1/sandbox_connect';

export class ApiClient {
  private transport: Transport;

  public sandboxService: PromiseClient<typeof SandboxService>;

  // Create a new ApiClient with the given token and baseUrl.
  // token can be either a Bearer token or an API token.
  constructor({
    baseUrl = 'https://api.jetpack.io',
    token,
  }: {
    baseUrl: string;
    token: string;
  }) {
    let tokenType = 'Bearer';
    // TODO: we should export some typeid types to make this check more robust.
    if (token?.startsWith('api_token')) {
      tokenType = 'Token';
    }
    this.transport = createConnectTransport({
      baseUrl,
      interceptors: [
        (next) => async (req) => {
          req.header.set('Authorization', `${tokenType} ${token}`);
          return await next(req);
        },
      ],
    });

    this.sandboxService = createPromiseClient(SandboxService, this.transport);
    // More services can be added here.
  }
}
