CREATE TABLE Option (
    id INT PRIMARY KEY,
    valor VARCHAR(255) NOT NULL
);

CREATE TABLE Proposal (
    id INT PRIMARY KEY,
    title VARCHAR(255) NOT NULL
);

CREATE TABLE Proposal_Option (
    id INT PRIMARY KEY,
    id_proposal INT NOT NULL,
    id_option INT NOT NULL,
    FOREIGN KEY (id_proposal) REFERENCES Proposal(id),
    FOREIGN KEY (id_option) REFERENCES Option(id),
    UNIQUE (id_proposal, id_option) 
);

CREATE TABLE Election (
    id INT PRIMARY KEY,    
    user_id INT NOT NULL,   
    proposal_id INT NOT NULL,
    option_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES User(id),
    FOREIGN KEY (proposal_id) REFERENCES Proposal(id),
    FOREIGN KEY (option_id) REFERENCES Option(id),
    UNIQUE (user_id, proposal_id),
    CONSTRAINT fk_proposal_option FOREIGN KEY (proposal_id, option_id)
    REFERENCES Proposal_Option(id_proposal, id_option)
);
