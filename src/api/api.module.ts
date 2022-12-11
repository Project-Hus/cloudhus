import { Module } from '@nestjs/common';
import { RouterModule } from '@nestjs/core';

import { PredApiModule } from './pred/api.pred.module';
import { UserApiModule } from './user/user.api.module';

@Module({
  imports: [
    RouterModule.register([
      {
        path: 'api',
        module: ApiModule,
        children: [
          {
            path: 'user',
            module: UserApiModule,
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
