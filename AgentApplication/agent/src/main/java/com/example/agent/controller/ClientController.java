package com.example.agent.controller;

import com.example.agent.model.Client;
import com.example.agent.model.Comment;
import com.example.agent.model.CompanyOwner;
import com.example.agent.model.JobPosition;
import com.example.agent.service.ClientService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(value = "/clients")
public class ClientController {

    @Autowired
    ClientService clientService;


    @PostMapping(path = "/create")
    public ResponseEntity<?> create(@RequestBody Client client) {
        return clientService.create(client);
    }


}
