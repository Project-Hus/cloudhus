/*
  Warnings:

  - Made the column `program_id` on table `routine` required. This step will fail if there are existing NULL values in that column.
  - Made the column `exercise_id` on table `routine` required. This step will fail if there are existing NULL values in that column.
  - Made the column `user_name` on table `user` required. This step will fail if there are existing NULL values in that column.
  - Made the column `password` on table `user` required. This step will fail if there are existing NULL values in that column.

*/
-- DropForeignKey
ALTER TABLE `routine` DROP FOREIGN KEY `Routine_exercise_id_fkey`;

-- DropForeignKey
ALTER TABLE `routine` DROP FOREIGN KEY `Routine_program_id_fkey`;

-- AlterTable
ALTER TABLE `routine` MODIFY `program_id` INTEGER NOT NULL,
    MODIFY `exercise_id` INTEGER NOT NULL;

-- AlterTable
ALTER TABLE `user` MODIFY `user_name` VARCHAR(191) NOT NULL,
    MODIFY `password` VARCHAR(191) NOT NULL;

-- AddForeignKey
ALTER TABLE `Routine` ADD CONSTRAINT `Routine_program_id_fkey` FOREIGN KEY (`program_id`) REFERENCES `Program`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `Routine` ADD CONSTRAINT `Routine_exercise_id_fkey` FOREIGN KEY (`exercise_id`) REFERENCES `Exercise`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
