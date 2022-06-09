import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { PostService } from 'src/app/service/post.service';

@Component({
  selector: 'app-create-post',
  templateUrl: './create-post.component.html',
  styleUrls: ['./create-post.component.css']
})
export class CreatePostComponent implements OnInit {

  constructor(private _postService: PostService, private router: Router) { }

  ngOnInit(): void {
  }

  post = {
    title: "",
    text: "",
    img: "",
    link: ""
  }
  errorMessage: any

  handleFileSelect(evt: any){
    var files = evt.target.files;
    var file = files[0];

    this.post.img = "data:image/"+ evt.target.files[0].name.split(".").pop() +";base64,"
  
    if (files && file) {
      var reader = new FileReader();
      //var extension = evt.target.files[0].name.split(".").pop()
      reader.onload =this._handleReaderLoaded.bind(this);
      reader.readAsBinaryString(file);
  }

}

_handleReaderLoaded(readerEvt : any) {
   var binaryString = readerEvt.target.result;
          this.post.img = ''.concat(this.post.img + btoa(binaryString));
          console.log(this.post.img)
  }

  createPost(){
    this._postService.post(this.post)
    .subscribe(
      error => this.errorMessage = <any>error);
  }
}

