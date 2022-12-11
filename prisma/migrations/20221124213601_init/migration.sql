-- CreateTable
CREATE TABLE `User` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `email_google` VARCHAR(191) NOT NULL,
    `token_google` VARCHAR(191) NOT NULL,
    `user_name` VARCHAR(191) NOT NULL,
    `password` VARCHAR(191) NOT NULL,
    `age` INTEGER NOT NULL,
    `sex` VARCHAR(191) NOT NULL,
    `height` DOUBLE NOT NULL,
    `arm_length` VARCHAR(191) NOT NULL,
    `leg_length` VARCHAR(191) NOT NULL,

    UNIQUE INDEX `User_email_google_key`(`email_google`),
    UNIQUE INDEX `User_user_name_key`(`user_name`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `Manager` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `manager_name` VARCHAR(191) NOT NULL,
    `password` VARCHAR(191) NOT NULL,

    UNIQUE INDEX `Manager_manager_name_key`(`manager_name`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

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
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
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

-- CreateTable
CREATE TABLE `Exercise` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `type_id` INTEGER NOT NULL,
    `name` VARCHAR(191) NOT NULL,
    `description` VARCHAR(191) NOT NULL,

    UNIQUE INDEX `Exercise_name_key`(`name`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

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
