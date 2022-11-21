import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { ApiModule } from './api/api.module';
import { PrismaService } from './services/prisma/prisma.service';
import { UserService } from './services/user/user.service';
import { AppService } from './app.service';


@Module({
  imports: [ ApiModule ],
  controllers: [ AppController ],
  providers: [ AppService, UserService, PrismaService ],
})
export class AppModule {}
