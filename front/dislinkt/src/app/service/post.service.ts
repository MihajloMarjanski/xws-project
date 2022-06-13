import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class PostService {

  constructor(private _http: HttpClient) { }

  private _baseUrl = 'https://localhost:8000/';
  private _postsForUser = this._baseUrl + 'post/user/';
  private _like = this._baseUrl + 'post/like';
  private _dislike = this._baseUrl + 'post/dislike';
  private _comment = this._baseUrl + 'post/comment';
  private _post = this._baseUrl + 'post';


  getPostsForUser(id: string|null|number) {
    if (id == null)
      id = 0
    return this._http.get<any>(this._postsForUser + id)
                           .pipe(tap(data =>  console.log('Iz service-a: ', data)),                         
                                catchError(this.handleError)); 
  }

  like(post: any): Observable<any>  {
    const body=JSON.stringify(post);
    return this._http.post(this._like, body)
  }

  dislike(post: any): Observable<any>  {
    const body=JSON.stringify(post);
    return this._http.post(this._dislike, body)
  }

  comment(comment: any): Observable<any>  {
    const body=JSON.stringify(comment);
    return this._http.post(this._comment, body)
  }

  post(post: any): Observable<any>  {
    const body=JSON.stringify(post);
    return this._http.post(this._post, body)
  }

  private handleError(err : HttpErrorResponse) {
    console.log(err.message);
    return Observable.throw(err.message);
    throw new Error('Method not implemented.');
  } 
}
