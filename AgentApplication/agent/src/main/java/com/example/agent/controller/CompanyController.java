package com.example.agent.controller;

import com.example.agent.model.CompanyOwner;
import com.example.agent.service.CompanyOwnerService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(value = "/company")
public class CompanyController {

    @Autowired
    private CompanyOwnerService companyOwnerService;


    @PostMapping(path = "/createOwner")
    public ResponseEntity<?> createCompanyOwner(@RequestBody CompanyOwner companyOwner) {
        companyOwnerService.save(companyOwner);
        return new ResponseEntity<>(HttpStatus.OK);
    }
}
