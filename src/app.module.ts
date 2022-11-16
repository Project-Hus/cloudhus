import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { ApiModule } from './api/api.module';
import { UserService } from './user/user.service';
import { PrismaService } from './services/prisma/prisma.service';


@Module({
  imports: [ ApiModule ],
  controllers: [ AppController ],
  providers: [ AppService, UserService, PrismaService ],
})
export class AppModule {}
