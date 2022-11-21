import { 
  Controller,
  Get,
  Param,
  Post,
  Body,
  Put,
  Delete
 } from '@nestjs/common';

import { User as UserModel } from '@prisma/client';
import { AppService } from './app.service';
import { UserService } from './services/user/user.service';

@Controller()
export class AppController {
  constructor(
    private readonly userService: UserService,
    private readonly appService: AppService
    ) {}
  
  @Get()
  sayHello (): string {
    return this.appService.getHello();
  }

  @Get('user/:id')
  async getUserById(@Param('id') id: number): Promise<UserModel> {
    return this.userService.user({ id: Number(id) });
  }
}
