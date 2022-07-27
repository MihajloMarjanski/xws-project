import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { JobOffer } from 'src/app/model/job-offer';
import { UserService } from 'src/app/service/user.service';

@Component({
  selector: 'app-create-job-offer',
  templateUrl: './create-job-offer.component.html',
  styleUrls: ['./create-job-offer.component.css']
})
export class CreateJobOfferComponent implements OnInit {

  errorMessage : string  = '';
  offer: JobOffer = {
    qualifications: '',
    companyName: '',
    jobInfo: '',
    jobPosition: '',
    apiKey: ''
  }

  constructor(public _userService: UserService, private router: Router, private _snackBar: MatSnackBar) { }

  ngOnInit(): void {
  }


  createOffer() {
    this._userService.getApiKey(localStorage.getItem("username")) 
        .subscribe(data => {
          this.offer.apiKey = data.apiKey
          
          this._userService.createOffer(this.offer)
              .subscribe(data => {
                this.router.navigateByUrl('/').then(() => {
                  this._snackBar.open('Job offer has been successfully created.', 'Close', {duration: 5000});
                });
                console.log('Dobio: ', data)},
              error => this.errorMessage = <any>error);  

          console.log('Dobio: ', data)},
        error => {
          this.errorMessage = <any>error
        });

  }
}
