namespace scheduler_service.DataAccess.Entities;

public class Timeslot
{
    public long Id { get; set; }

    public long ApartmentId { get; set; }

    public DateTime StartDatetime { get; set; }

    public DateTime EndDatetime { get; set; }

    public virtual IEnumerable<Reservation> Reservations { get; set; } = [];
}