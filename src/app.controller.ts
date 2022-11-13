import { 
  Controller,
  Get,
  Param,
  Post,
  Body,
  Put,
  Delete
 } from '@nestjs/common';

import { AppService } from './app.service';
import { UserService } from './user/user.service';

import { User as UserModel } from '@prisma/client';

@Controller()
export class AppController {
  constructor(
    private readonly appService: AppService,
    private readonly userService: UserService,
    ) {}

  @Get('user/:id')
  async getUserById(@Param('id') id: number): Promise<UserModel> {
    return this.userService.user({ id: Number(id) });
  }
}
