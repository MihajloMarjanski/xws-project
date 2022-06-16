import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PasswrodlessComponent } from './passwrodless.component';

describe('PasswrodlessComponent', () => {
  let component: PasswrodlessComponent;
  let fixture: ComponentFixture<PasswrodlessComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PasswrodlessComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PasswrodlessComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
