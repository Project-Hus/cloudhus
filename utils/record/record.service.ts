import { Injectable } from '@nestjs/common';
import { RecordFixed } from 'src/dto/RecordFixed';
import { RecordInput } from 'src/dto/RecordInput';
import { RecordWeekly } from 'src/dto/RecordWeekly';

@Injectable()
export class RecordService {
  /**
   * It attaches the fixed part and variable part of the model input.
   * @param recordFixed Constant values like sex, arm length etc
   * @param recordWeekly Variables like bodyweight etc
   * @returns 
   */
  processRecords(recordFixed:RecordFixed, recordWeekly:RecordWeekly[]): RecordInput[] {
    return [
        {test:'hi', a:3, b:5},
        {test:'no', a:1, b:4},
        {test:'why', a:2, b:7}
    ]
  }
}
