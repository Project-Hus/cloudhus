import { Module } from '@nestjs/common';
import { ApiController } from './api.controller';
import { PredService } from './services/pred.service';

@Module({
  imports: [],
  controllers: [ ApiController ],
  providers: [ PredService ],
})
export class ApiModule {}
