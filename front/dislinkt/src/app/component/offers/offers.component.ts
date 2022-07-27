import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { JobOffer } from 'src/app/model/job-offer';
import { UserService } from 'src/app/service/user.service';

@Component({
  selector: 'app-offers',
  templateUrl: './offers.component.html',
  styleUrls: ['./offers.component.css']
})
export class OffersComponent implements OnInit {
  searchField: string = '';
  displayedColumns: string[] = ['companyName', 'jobPosition', 'jobInfo'];
  errorMessage : string  = '';
  show : boolean = false
  offers: JobOffer[] = []
  offer: JobOffer = {
    companyName: "",
    apiKey: "",
    jobInfo:  "",
    jobPosition: "",
    qualifications: '',
  }

  constructor(public _userService: UserService, private _snackBar: MatSnackBar, private router: Router) { }

  ngOnInit(): void {
    this.searchTable()
  }



  showOffer(selectedOffer: JobOffer) {
    this.offer = selectedOffer
    this.show = true
  }

  searchTable() {
    this._userService.searchOffers(this.searchField)
          .subscribe(data => {
              this.offers = data.offers
              this.show = false
            console.log('Dobio: ', data)},
          error => this.errorMessage = <any>error);  
    
  }

}
