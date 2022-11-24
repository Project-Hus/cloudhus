/*
  Warnings:

  - You are about to drop the column `user_name` on the `manager` table. All the data in the column will be lost.
  - You are about to drop the `bodyinfo` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `program` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `programrec` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `routine` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `routinerec` table. If the table is not empty, all the data it contains will be lost.
  - A unique constraint covering the columns `[name]` on the table `Exercise` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[manager_name]` on the table `Manager` will be added. If there are existing duplicate values, this will fail.
  - Added the required column `type_id` to the `Exercise` table without a default value. This is not possible if the table is not empty.
  - Added the required column `manager_name` to the `Manager` table without a default value. This is not possible if the table is not empty.
  - Added the required column `updated_at` to the `Manager` table without a default value. This is not possible if the table is not empty.
  - Added the required column `arm_length` to the `User` table without a default value. This is not possible if the table is not empty.
  - Added the required column `height` to the `User` table without a default value. This is not possible if the table is not empty.
  - Added the required column `leg_length` to the `User` table without a default value. This is not possible if the table is not empty.
  - Added the required column `updated_at` to the `User` table without a default value. This is not possible if the table is not empty.

*/
-- DropForeignKey
ALTER TABLE `bodyinfo` DROP FOREIGN KEY `BodyInfo_user_id_fkey`;

-- DropForeignKey
ALTER TABLE `programrec` DROP FOREIGN KEY `ProgramRec_program_id_fkey`;

-- DropForeignKey
ALTER TABLE `routine` DROP FOREIGN KEY `Routine_exercise_id_fkey`;

-- DropForeignKey
ALTER TABLE `routine` DROP FOREIGN KEY `Routine_program_id_fkey`;

-- DropForeignKey
ALTER TABLE `routinerec` DROP FOREIGN KEY `RoutineRec_programrec_id_fkey`;

-- DropForeignKey
ALTER TABLE `routinerec` DROP FOREIGN KEY `RoutineRec_routine_id_fkey`;

-- AlterTable
ALTER TABLE `exercise` ADD COLUMN `type_id` INTEGER NOT NULL;

-- AlterTable
ALTER TABLE `manager` DROP COLUMN `user_name`,
    ADD COLUMN `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    ADD COLUMN `manager_name` VARCHAR(191) NOT NULL,
    ADD COLUMN `updated_at` DATETIME(3) NOT NULL;

-- AlterTable
ALTER TABLE `user` ADD COLUMN `arm_length` VARCHAR(191) NOT NULL,
    ADD COLUMN `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    ADD COLUMN `height` DOUBLE NOT NULL,
    ADD COLUMN `leg_length` VARCHAR(191) NOT NULL,
    ADD COLUMN `updated_at` DATETIME(3) NOT NULL;

-- DropTable
DROP TABLE `bodyinfo`;

-- DropTable
DROP TABLE `program`;

-- DropTable
DROP TABLE `programrec`;

-- DropTable
DROP TABLE `routine`;

-- DropTable
DROP TABLE `routinerec`;

-- CreateTable
CREATE TABLE `TrainingProgramType` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `type` VARCHAR(191) NOT NULL,

    UNIQUE INDEX `TrainingProgramType_type_key`(`type`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `TrainingProgram` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,
    `type_id` INTEGER NOT NULL,
    `author` INTEGER NOT NULL,
    `name` VARCHAR(191) NOT NULL,
    `description` VARCHAR(191) NOT NULL,
    `vector` INTEGER NULL,

    UNIQUE INDEX `TrainingProgram_name_key`(`name`),
    UNIQUE INDEX `TrainingProgram_vector_key`(`vector`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `TrainingProgramRec` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `program_id` INTEGER NOT NULL,
    `user_id` INTEGER NOT NULL,
    `start` DATETIME(3) NOT NULL,
    `end` DATETIME(3) NOT NULL,
    `comment` VARCHAR(191) NOT NULL,
    `weight` DOUBLE NOT NULL,
    `fat_rate` DOUBLE NOT NULL,
    `squat` DOUBLE NOT NULL,
    `benchpress` DOUBLE NOT NULL,
    `deadlift` DOUBLE NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ProgramVector` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `c0` INTEGER NOT NULL,
    `c1` INTEGER NOT NULL,
    `c2` INTEGER NOT NULL,
    `c3` INTEGER NOT NULL,
    `c4` INTEGER NOT NULL,
    `c5` INTEGER NULL,
    `c6` INTEGER NULL,
    `c7` INTEGER NULL,
    `c8` INTEGER NULL,
    `c9` INTEGER NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `WeekRoutine` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `program_id` INTEGER NOT NULL,
    `description` VARCHAR(191) NOT NULL,
    `order` INTEGER NOT NULL,

    UNIQUE INDEX `WeekRoutine_order_key`(`order`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `WeekRoutineRec` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `week_routine_id` INTEGER NOT NULL,
    `program_id` INTEGER NOT NULL,
    `comment` VARCHAR(191) NOT NULL,
    `weight` DOUBLE NOT NULL,
    `fat_rate` DOUBLE NOT NULL,
    `squat` DOUBLE NOT NULL,
    `benchpress` DOUBLE NOT NULL,
    `deadlift` DOUBLE NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `DayRoutine` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `week_routine_id` INTEGER NOT NULL,
    `order` INTEGER NOT NULL,
    `exercise_id` INTEGER NOT NULL,
    `reps` INTEGER NOT NULL,
    `description` VARCHAR(191) NOT NULL,

    UNIQUE INDEX `DayRoutine_order_key`(`order`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `DayRoutineRec` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `week_routine_id` INTEGER NOT NULL,
    `the_day` DATETIME(3) NOT NULL,
    `reps` INTEGER NOT NULL,
    `comment` VARCHAR(191) NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ExerciseType` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `type` VARCHAR(191) NOT NULL,
    `description` VARCHAR(191) NOT NULL,

    UNIQUE INDEX `ExerciseType_type_key`(`type`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateIndex
CREATE UNIQUE INDEX `Exercise_name_key` ON `Exercise`(`name`);

-- CreateIndex
CREATE UNIQUE INDEX `Manager_manager_name_key` ON `Manager`(`manager_name`);

-- AddForeignKey
ALTER TABLE `TrainingProgram` ADD CONSTRAINT `TrainingProgram_type_id_fkey` FOREIGN KEY (`type_id`) REFERENCES `TrainingProgramType`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `TrainingProgram` ADD CONSTRAINT `TrainingProgram_author_fkey` FOREIGN KEY (`author`) REFERENCES `User`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `TrainingProgram` ADD CONSTRAINT `TrainingProgram_vector_fkey` FOREIGN KEY (`vector`) REFERENCES `ProgramVector`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `TrainingProgramRec` ADD CONSTRAINT `TrainingProgramRec_program_id_fkey` FOREIGN KEY (`program_id`) REFERENCES `TrainingProgram`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `TrainingProgramRec` ADD CONSTRAINT `TrainingProgramRec_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `User`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `WeekRoutine` ADD CONSTRAINT `WeekRoutine_program_id_fkey` FOREIGN KEY (`program_id`) REFERENCES `TrainingProgram`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `WeekRoutineRec` ADD CONSTRAINT `WeekRoutineRec_week_routine_id_fkey` FOREIGN KEY (`week_routine_id`) REFERENCES `WeekRoutine`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `WeekRoutineRec` ADD CONSTRAINT `WeekRoutineRec_program_id_fkey` FOREIGN KEY (`program_id`) REFERENCES `TrainingProgramRec`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `DayRoutine` ADD CONSTRAINT `DayRoutine_week_routine_id_fkey` FOREIGN KEY (`week_routine_id`) REFERENCES `WeekRoutine`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `DayRoutine` ADD CONSTRAINT `DayRoutine_exercise_id_fkey` FOREIGN KEY (`exercise_id`) REFERENCES `Exercise`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `DayRoutineRec` ADD CONSTRAINT `DayRoutineRec_week_routine_id_fkey` FOREIGN KEY (`week_routine_id`) REFERENCES `WeekRoutineRec`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `Exercise` ADD CONSTRAINT `Exercise_type_id_fkey` FOREIGN KEY (`type_id`) REFERENCES `ExerciseType`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
