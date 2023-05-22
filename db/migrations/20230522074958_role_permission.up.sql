CREATE TABLE Roles (
  id uuid PRIMARY KEY,
  name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE Permissions (
  id uuid PRIMARY KEY,
  name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE UserRoles (
  user_id uuid REFERENCES Users(id),
  role_id uuid REFERENCES Roles(id),
  PRIMARY KEY (user_id, role_id)
);

CREATE TABLE RolePermissions (
  role_id uuid REFERENCES Roles(id),
  permission_id uuid REFERENCES Permissions(id),
  PRIMARY KEY (role_id, permission_id)
);
