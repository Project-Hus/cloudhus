import { Module } from '@nestjs/common';
import { ApiModule } from './api/api.module';

import { PrismaService } from './prisma_services/prisma/prisma.service';

@Module({
  imports: [ ApiModule ],
  controllers: [  ],
  providers: [ PrismaService],
})
export class AppModule {}
