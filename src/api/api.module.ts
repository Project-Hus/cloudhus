import { Module } from '@nestjs/common';
import { RouterModule } from '@nestjs/core';

import { PredApiModule } from './pred/api.pred.module';
import { AuthApiModule } from './auth/auth.api.module';

@Module({
  imports: [
    RouterModule.register([
      {
        path: 'api',
        module: ApiModule,
        children: [
          {
            path: 'auth',
            module: AuthApiModule,
          },
          {
            path: 'pred',
            module: PredApiModule,
          },
        ],
      },
    ]),
  ],
})
export class ApiModule {}
