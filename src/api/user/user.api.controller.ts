import { 
    Controller,
    Get,
    Param,
    Post,
    Body,
    Put,
    Delete
   } from '@nestjs/common';
import { UserRegisterForm } from 'src/dto/UserRegisterForm';
import { UserPrismaService } from 'src/prisma_services/user/user.prisma.service';
import { UserApiService } from './services/user.api.service';

  @Controller()
  export class UserApiController {
    constructor(
      private readonly userPrismaService : UserPrismaService,

      private readonly userApiService : UserApiService
      ) {}

    /**
     * API for register
     * @param records Training records for 24 weeks
     * @returns Best 3 methods and the results
     */
    @Post('register') // POST /api/register
    async getPred(
      @Body('user_register_form') user_register_form: UserRegisterForm,
      ): Promise<string> {
      return this.userApiService.register(user_register_form);
    } 
  }
  