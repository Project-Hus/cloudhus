import { INestApplication, Injectable, OnModuleInit } from "@nestjs/common";
import { PrismaClient } from '@prisma/client';

// onModuleInit is optional - if you leave it out, Prisma will connect lazily on its first call to the DB.
@Injectable()
export class PrismaService extends PrismaClient implements OnModuleInit {
    async onModuleInit() {
        await this.$connect();
    }

    async enableShutdownHooks(app: INestApplication) {
        this.$on('beforeExit', async ()=> {
            await app.close();
        });
    }
}