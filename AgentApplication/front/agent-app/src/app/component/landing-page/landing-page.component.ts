import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-landing-page',
  templateUrl: './landing-page.component.html',
  styleUrls: ['./landing-page.component.css']
})
export class LandingPageComponent implements OnInit {

  slides = [
    {
      image: /* "../images/slide1.jpg" */
        "https://images.pexels.com/photos/629167/pexels-photo-629167.jpeg?auto=compress&cs=tinysrgb&dpr=2&h=650&w=940"
    },
    {
      image:
        "https://images.pexels.com/photos/294674/pexels-photo-294674.jpeg?auto=compress&cs=tinysrgb&dpr=2&h=650&w=940"
    },
    {
      image:
        "https://images.pexels.com/photos/1105386/pexels-photo-1105386.jpeg?auto=compress&cs=tinysrgb&dpr=2&h=650&w=940"
    },
    {
      image:
        "https://images.pexels.com/photos/2582868/pexels-photo-2582868.jpeg?auto=compress&cs=tinysrgb&dpr=2&h=650&w=940"
    },
    {
      image:
        "https://images.pexels.com/photos/1612351/pexels-photo-1612351.jpeg?auto=compress&cs=tinysrgb&dpr=3&h=750&w=1260"
    }
  ];

  constructor() { }

  ngOnInit(): void {
  }

}
