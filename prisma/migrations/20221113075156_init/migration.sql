-- CreateTable
CREATE TABLE `User` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `email_google` VARCHAR(191) NOT NULL,
    `token_google` VARCHAR(191) NOT NULL,
    `user_name` VARCHAR(191) NULL,
    `password` VARCHAR(191) NULL,
    `age` INTEGER NOT NULL,
    `sex` VARCHAR(191) NOT NULL,

    UNIQUE INDEX `User_email_google_key`(`email_google`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `BodyInfo` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `user_id` INTEGER NULL,
    `date` DATETIME(3) NOT NULL,
    `height` DOUBLE NOT NULL,
    `weight` DOUBLE NOT NULL,
    `fat_rate` DOUBLE NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `Manager` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `user_name` VARCHAR(191) NOT NULL,
    `password` VARCHAR(191) NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `Program` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(191) NOT NULL,
    `description` VARCHAR(191) NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `Routine` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `program_id` INTEGER NULL,
    `week` INTEGER NOT NULL,
    `day` INTEGER NOT NULL,
    `order` INTEGER NOT NULL,
    `exercise_id` INTEGER NULL,
    `reps` INTEGER NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `Exercise` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(191) NOT NULL,
    `description` VARCHAR(191) NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- AddForeignKey
ALTER TABLE `BodyInfo` ADD CONSTRAINT `BodyInfo_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `User`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `Routine` ADD CONSTRAINT `Routine_program_id_fkey` FOREIGN KEY (`program_id`) REFERENCES `Program`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `Routine` ADD CONSTRAINT `Routine_exercise_id_fkey` FOREIGN KEY (`exercise_id`) REFERENCES `Exercise`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;
