import { Injectable } from '@nestjs/common';
import { RecordFixed } from 'src/dto/RecordFixed';
import { RecordWeekly } from 'src/dto/RecordWeekly';

@Injectable()
export class RecordService {
  /**
   * It attaches the fixed part and variable part of the model input.
   * @param recordFixed Constant values like sex, arm length etc
   * @param recordWeekly Variables like bodyweight etc
   * @returns 
   */
  processRecords(recordFixed:RecordFixed, recordWeekly:RecordWeekly[]): [] {
    return []
  }
}
