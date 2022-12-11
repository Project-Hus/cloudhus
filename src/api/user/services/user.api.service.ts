import { Injectable } from '@nestjs/common';

/* DTOs */
import { UserRegisterForm } from 'src/dto/UserRegisterForm';

/* Services */
import { UserPrismaService } from 'src/prisma_services/user/user.prisma.service';

/**
 * API Service for User processing
 */
@Injectable()
export class UserApiService {
  constructor(
    private readonly userPrismaService : UserPrismaService,
  ) {}

  async register(user:UserRegisterForm): Promise<any> {
    try{
  
    } catch (error) {
    }
  }
}
