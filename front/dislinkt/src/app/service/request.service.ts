import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
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
  private _getRequests = this._baseUrl + 'requests/getAll/'
  private _acceptRequest = this._baseUrl + 'requests/acceptRequest/'
  private _declineRequest = this._baseUrl + 'requests/declineRequest/'
  
  
  constructor(private _http: HttpClient) { }
  

  declineRequest(senderId: number, receiverId: number) : Observable<any> {
    return this._http.put<any>(this._declineRequest + senderId + "/" + receiverId, {})
                  .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                  catchError(this.handleError)); 
  }

  acceptRequest(senderId: number, receiverId: number) : Observable<any> {
    return this._http.put<any>(this._acceptRequest + senderId + "/" + receiverId, {})
                  .pipe(tap(data =>  console.log('All: ' + JSON.stringify(data))),
                  catchError(this.handleError)); 
  }

  getRequests(logedUserId: string | null) {
    return this._http.get<any>(this._getRequests + logedUserId)
            .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                catchError(this.handleError));
  }

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
