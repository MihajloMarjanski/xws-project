import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { PostService } from 'src/app/service/post.service';

@Component({
  selector: 'app-view-post',
  templateUrl: './view-post.component.html',
  styleUrls: ['./view-post.component.css']
})
export class ViewPostComponent implements OnInit {

  constructor(private _postService: PostService, private router: Router) { }

  post: any
  likeCount: number = 0
  dislikeCount: number = 0
  postToLike = {
    post: ""
  }
  comment = {
    post: "",
    text: ""
  }
  errorMessage :any

  ngOnInit(): void {
    this.post = history.state.data
    this.likeCount = this.post.like.length
    this.dislikeCount = this.post.dislike.length
    this.postToLike.post = this.post.id
    this.comment.post = this.post.id
  }

  like(){
    this._postService.like(this.postToLike)
        .subscribe(
          error => this.errorMessage = <any>error);
  }

  dislike(){
    this._postService.dislike(this.postToLike)
        .subscribe(
          error => this.errorMessage = <any>error);
  }
  
  makeComment(){
    this._postService.comment(this.comment)
    .subscribe(
      error => this.errorMessage = <any>error);
    this.comment.text = ""
  }

}
