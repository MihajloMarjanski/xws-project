package com.example.agent.model.dto;

import com.example.agent.model.Company;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import javax.persistence.Column;
import java.util.Date;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class OwnerWithCompany {
    private Integer id;
    private String firstName;
    private String lastName;
    private String email;
    private String username;
    private String password;
    private boolean isActivated = true;
    private Integer forgotten;
    private Integer pin;
    private String salt;
    private Integer missedPasswordCounter;
    private boolean isBlocked;
    private Date blockedDate;

    private Company company;
}
