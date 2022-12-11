import { Injectable } from '@nestjs/common';
import { UserSignInDTO } from 'src/dto/UserSignIn.dto';

/* DTOs */
import { UserSignUpDTO } from 'src/dto/UserSignUp.dto';

/* Services */
import { UserPrismaService } from 'src/prisma_services/user/user.prisma.service';

/**
 * API Service for User processing
 */
@Injectable()
export class AuthApiService {
  constructor(
    private readonly userPrismaService : UserPrismaService,
  ) {}

  async sign_up(user_sign_up_dto:UserSignUpDTO): Promise<boolean> {
    try{
      this.userPrismaService
    } catch (error) {
    }
  }

  async sign_in(user_sign_in_dto: UserSignInDTO): Promise<boolean> {
    try{
  
    } catch (error) {
    }
  }

  async sign_out(): Promise<boolean> {
    try{
  
    } catch (error) {
    }
  }

  async delete_account(): Promise<boolean> {
    try{
  
    } catch (error) {
    }
  }
}
