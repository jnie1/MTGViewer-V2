package container

import (
	"github.com/jnie1/MTGViewer-V2/database"
)

func AddContainer(cont Container) error {
	db := database.Instance()

	_, err := db.Exec(`
	INSERT INTO container (container_name, capacity, deletion_mark) 
	VALUES ($1, $2, FALSE)`, cont.Name, cont.Capacity)

	return err
}

func GetContainer(containerId int) (Container, error) {
	container := Container{}

	db := database.Instance()

	row := db.QueryRow(`
	SELECT container, capacity
	FROM container
	WHERE container_id = $1`, containerId)

	err := row.Scan(&container.Name, &container.Capacity, &container.MarkForDeletion)

	return container, err
}

func UpdateContainerName(containerId int, newName string) (Container, error) {

	db := database.Instance()

	row := db.QueryRow(`
	UPDATE container
	SET container_name = $1
	WHERE container_id = $2`, newName, containerId)

	updatedContainer := Container{}
	err := row.Scan(&updatedContainer.Name, &updatedContainer.Capacity, &updatedContainer.MarkForDeletion)

	if err != nil {
		return Container{}, err
	}

	return updatedContainer, nil
}

func UpdateContainerCapacity(containerId int, newCapacity int) (Container, error) {

	db := database.Instance()

	row := db.QueryRow(`
	UPDATE container
	SET capacity = $1
	WHERE container_id = $2`, newCapacity, containerId)

	updatedContainer := Container{}
	err := row.Scan(&updatedContainer.Name, &updatedContainer.Capacity, &updatedContainer.MarkForDeletion)

	if err != nil {
		return Container{}, err
	}

	return updatedContainer, nil
}

func UpdateContainerDeletionMark(containerId int, delMark bool) (Container, error) {

	db := database.Instance()

	row := db.QueryRow(`
	UPDATE container
	SET deletion_mark = $1
	WHERE container_id = $2`, delMark, containerId)

	updatedContainer := Container{}
	err := row.Scan(&updatedContainer.Name, &updatedContainer.Capacity, &updatedContainer.MarkForDeletion)

	if err != nil {
		return Container{}, err
	}

	return updatedContainer, nil
}

func DeleteContainer(containerId int) error {

	db := database.Instance()

	_, err := db.Exec(`
	DELETE FROM container
	WHERE container_id = $1`, containerId)

	return err
}
