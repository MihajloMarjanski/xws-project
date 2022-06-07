import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Credentials } from '../model/credentials';
import { Observable, throwError } from 'rxjs';
import { map, catchError, tap } from 'rxjs/operators';
import { User } from '../model/user';
import { Experience } from '../model/experience';
import { Interest } from '../model/interest';
import { JobOffer } from '../model/job-offer';
import { Message } from '../model/message';

@Injectable({
  providedIn: 'root'
})
export class RequestService {
  
  private _baseUrl = 'https://localhost:8000/';
  private _sendConnectRequest = this._baseUrl + 'requests/sendRequest/'
  private _areConnected = this._baseUrl + 'connection/'
  private _findConnections = this._baseUrl + 'connections/'
  private _sendMessage = this._baseUrl + 'message/send/'
  private _findMessages = this._baseUrl + 'messages/'
  private _getNotifications = this._baseUrl + 'notifications/'
  
  constructor(private _http: HttpClient) { }
  

  getNotifications(logedUserId: string | null) {
    return this._http.get<any>(this._getNotifications + logedUserId)
            .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                catchError(this.handleError));
  }

  sendMessage(newMessage: Message) : Observable<any>  {
    const body=JSON.stringify(newMessage);
    return this._http.post(this._sendMessage + newMessage.senderId + '/' + newMessage.receiverId, body)
  }

  getMessagesForConnection(logedUserId: number, id: number) {
    return this._http.get<any>(this._findMessages + logedUserId + '/' + id)
            .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                catchError(this.handleError));
  }

  findConnections(logedUserId: number) {
    return this._http.get<any>(this._findConnections + logedUserId)
            .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                catchError(this.handleError));
  }

  sendConnectRequest(logedUserId: string | null, wantedUserId: number) : Observable<any> {
    return this._http.put<any>(this._sendConnectRequest + logedUserId + "/" + wantedUserId, {})
                  .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                  catchError(this.handleError)); 
  }
  
  areConnected(id1: string | null, id2: number) {
    return this._http.get<any>(this._areConnected + id1 + "/" + id2)
            .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                catchError(this.handleError));
  }
  
  private handleError(err : HttpErrorResponse) {
    console.log(err.message);
    return Observable.throw(err.message);
    throw new Error('Method not implemented.');
  } 
}
