import { 
    Controller,
    Post,
    Body,
    Delete
   } from '@nestjs/common';

/* Services */
import { AuthApiService } from './services/auth.api.service';

/* DTOs */
import { UserSignInDTO } from 'src/dto/UserSignIn.dto';
import { UserSignUpDTO } from 'src/dto/UserSignUp.dto';

  @Controller()
  export class AuthApiController {
    constructor(
      private readonly authApiService : AuthApiService
      ) {}

    /**
     * POST /api/auth/signup
     * API for signing up 
     * @param user_sign_up_dto user sign up form
     * @returns True/ False (succeeded/failed)
     */
    @Post('signup')
    async sign_up(
      @Body('user_sign_up_form_dto') user_sign_up_dto: UserSignUpDTO,
      ): Promise<boolean> {
      return this.authApiService.sign_up(user_sign_up_dto);
    }

    /**
     * POST /api/auth/signin
     * API for signing in
     * @param user_register_in_dto user sign in form
     * @returns True/ False (succeeded/failed)
     */
     @Post('signin')
     async sign_in(
       @Body('user_sign_in_dto') user_sign_in_dto: UserSignInDTO,
       ): Promise<boolean> {
       return this.authApiService.sign_in(user_sign_in_dto);
     }

     /**
     * POST /api/auth/signup
     * API for signing out
     * @returns True/ False (succeeded/failed)
     */
    @Post('signout')
    async sign_out(): Promise<boolean> {
      return this.authApiService.sign_out();
    }

    /**
     * API for signing up 
     * @param user_register_form user register form
     * @returns True/ False (succeeded/failed)
     */
     @Delete('delete') // POST /api/auth/signup
     async delete_account(
       ): Promise<boolean> {
       return this.authApiService.delete_account();
     }
  }
  