import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { PostService } from 'src/app/service/post.service';

@Component({
  selector: 'app-posts',
  templateUrl: './posts.component.html',
  styleUrls: ['./posts.component.css']
})
export class PostsComponent implements OnInit {

  constructor(private _postService: PostService, private router: Router) { }
  posts: any[] = []
  errorMessage: any
  displayedColumns: string[] = ['title'];
  
  

  ngOnInit(): void {
    this.getPostsForUser(history.state.data)

  }

  viewPost(post: any){  
    this.router.navigate(["/view-post"],{state: {data: post}})
  }

  getPostsForUser(id: string|null|number){
    this._postService.getPostsForUser(id)
    .subscribe(data => {
        this.posts = data.post
      console.log('Dobio: ', data)},
    error => this.errorMessage = <any>error);  
  }
}
